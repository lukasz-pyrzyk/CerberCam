package main

import (
	"archive/zip"
	"bufio"
	"bytes"
	"fmt"
	"image"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/tensorflow/tensorflow/tensorflow/go"
	"github.com/tensorflow/tensorflow/tensorflow/go/op"
)

func recognize(msg Message, optionString string) {
	options := new(tensorflow.SessionOptions)
	options.Target = optionString

	log.Info("Recognizing image...")
	log.Info("Looking for model file...")

	// Load the serialized GraphDef from a file.
	modelfile, labelsfile, err := modelFiles(*modeldir)
	failOnError(err, "No model file in model directory")

	log.Info("Loading model file...")

	model, err := ioutil.ReadFile(modelfile)
	failOnError(err, "Unable to load model file!")

	log.Info("Constructing graph...")

	// Construct an in-memory graph from the serialized form.
	graph := tensorflow.NewGraph()
	if err := graph.Import(model, ""); err != nil {
		log.Fatal(err)
	}

	log.Info("Starting new session...")

	// Create a session for inference over graph.
	session, err := tensorflow.NewSession(graph, options)
	failOnError(err, "Unable to start new tensorflow session")
	defer session.Close()

	log.Info("Creating tensor from image...")

	// Run inference on thestImageFilename.
	// For multiple images, session.Run() can be called in a loop (and
	// concurrently). Furthermore, images can be batched together since the
	// model accepts batches of image data as input.
	tensor, err := makeTensorFromImage(msg.Photo)
	failOnError(err, "Unable to make tensor from image!")

	log.Info("Running recognition...")

	output, err := session.Run(
		map[tensorflow.Output]*tensorflow.Tensor{
			graph.Operation("input").Output(0): tensor,
		},
		[]tensorflow.Output{
			graph.Operation("output").Output(0),
		},
		nil)

	failOnError(err, "Error while running a tensorflow session")

	log.Info("Gathering output...")

	// output[0].Value() is a vector containing probabilities of
	// labels for each image in the "batch". The batch size was 1.
	// Find the most probably label index.
	probabilities := output[0].Value().([][]float32)[0]
	printBestLabel(probabilities, labelsfile)
}

// Conver the image in filename to a Tensor suitable as input to the Inception model.
func makeTensorFromImage(data []byte) (*tensorflow.Tensor, error) {
	reader := bytes.NewReader(data)
	img, _, err := image.Decode(reader)
	if err != nil {
		return nil, err
	}

	// Represent the image as [H][W][B,G,R]byte
	contents := make([][][3]byte, img.Bounds().Size().Y)
	for y := 0; y < len(contents); y++ {
		contents[y] = make([][3]byte, img.Bounds().Size().X)
		for x := 0; x < len(contents[y]); x++ {
			px := x + img.Bounds().Min.X
			py := y + img.Bounds().Min.Y
			r, g, b, _ := img.At(px, py).RGBA()

			// image.Image uses 16-bits for each color.
			// We want 8-bits.
			contents[y][x][0] = byte(b >> 8)
			contents[y][x][1] = byte(g >> 8)
			contents[y][x][2] = byte(r >> 8)
		}
	}

	tensor, err := tensorflow.NewTensor(contents)
	if err != nil {
		return nil, err
	}

	// Construct a graph to normalize the image
	graph, input, output, err := constructGraphToNormalizeImage()
	if err != nil {
		return nil, err
	}

	// Execute that graph to normalize this one image
	session, err := tensorflow.NewSession(graph, nil)
	if err != nil {
		return nil, err
	}

	defer session.Close()

	normalized, err := session.Run(
		map[tensorflow.Output]*tensorflow.Tensor{input: tensor},
		[]tensorflow.Output{output},
		nil)

	if err != nil {
		return nil, err
	}

	return normalized[0], nil
}

// The inception model takes as input the image described by a Tensor in a very
// specific normalized format (a particular image size, shape of the input tensor,
// normalized pixel values etc.).
//
// This function constructs a graph of TensorFlow operations which takes as input
// the raw pixel values of an image in the form of a Tensor of shape [Height, Width, 3]
// and returns a tensor suitable for input to the inception model.
//
// T[y][x] is the (Blue, Green, Red) values of the pixel at position (x, y) in the image,
// with each color value represented as a single byte.
func constructGraphToNormalizeImage() (graph *tensorflow.Graph, input, output tensorflow.Output, err error) {
	// Some constants specific to the pre-trained model at:
	// https://storage.googleapis.com/download.tensorflow.org/models/inception5h.zip
	//
	// - The model was trained after with images scaled to 224x224 pixels.
	// - The colors, represented as R, G, B in 1-byte each were converted to
	//   float using (value - Mean)/Scale.
	//
	// If using a different pre-trained model, the values will have to be adjusted.
	const (
		H, W  = 224, 224
		Mean  = float32(117)
		Scale = float32(1)
	)
	// - input is a 3D tensor of shape [Height, Width, Colors=3], where
	//   each pixel is represented as a triplet of 1-byte colors
	// - ResizeBilinear (and the inception model) takes a 4D tensor of shape
	//   [BatchSize, Height, Width, Colors=3], where each pixel is
	//   represented as a triplet of floats
	// - Apply normalization on each pixel and use ExpandDims to make
	//   this single image be a "batch" of size 1 for ResizeBilinear.
	s := op.NewScope()

	input = op.Placeholder(s, tensorflow.Uint8)
	output = op.Div(s,
		op.Sub(s,
			op.ResizeBilinear(s,
				op.ExpandDims(s,
					op.Cast(s, input, tensorflow.Float),
					op.Const(s.SubScope("make_batch"), int32(0))),
				op.Const(s.SubScope("size"), []int32{H, W})),
			op.Const(s.SubScope("mean"), Mean)),
		op.Const(s.SubScope("scale"), Scale))
	graph, err = s.Finalize()
	return graph, input, output, err
}

func printBestLabel(probabilities []float32, labelsFile string) {
	bestIdx := 0
	for i, p := range probabilities {
		if p > probabilities[bestIdx] {
			bestIdx = i
		}
	}
	// Found a best match, now read the string from the labelsFile where
	// there is one line per label.
	file, err := os.Open(labelsFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var labels []string
	for scanner.Scan() {
		labels = append(labels, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Errorf("ERROR: failed to read %s: %v", labelsFile, err)
	}
	fmt.Printf("BEST MATCH: (%2.0f%% likely) %s\n", probabilities[bestIdx]*100.0, labels[bestIdx])
}

// TODO: refactor those methods
func modelFiles(dir string) (modelfile, labelsfile string, err error) {
	const URL = "https://storage.googleapis.com/download.tensorflow.org/models/inception5h.zip"

	var (
		model   = filepath.Join(dir, "tensorflow_inception_graph.pb")
		labels  = filepath.Join(dir, "imagenet_comp_graph_label_strings.txt")
		zipfile = filepath.Join(dir, "inception5h.zip")
	)

	if filesExist(model, labels) == nil {
		return model, labels, nil
	}

	log.Warningf("Did not find model in '%s' downloading from '%s'", dir, URL)

	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", "", err
	}

	if err := download(URL, zipfile); err != nil {
		return "", "", fmt.Errorf("failed to download %v - %v", URL, err)
	}

	if err := unzip(dir, zipfile); err != nil {
		return "", "", fmt.Errorf("failed to extract contents from model archive: %v", err)
	}

	os.Remove(zipfile)
	return model, labels, filesExist(model, labels)
}

func filesExist(files ...string) error {
	for _, f := range files {
		if _, err := os.Stat(f); err != nil {
			return fmt.Errorf("unable to stat %s: %v", f, err)
		}
	}
	return nil
}

func download(URL, filename string) error {
	resp, err := http.Get(URL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = io.Copy(file, resp.Body)
	return err
}

func unzip(dir, zipfile string) error {
	r, err := zip.OpenReader(zipfile)
	if err != nil {
		return err
	}
	defer r.Close()
	for _, f := range r.File {
		src, err := f.Open()
		if err != nil {
			return err
		}
		log.Infof("Extracting from archive '%s'...", f.Name)
		dst, err := os.OpenFile(filepath.Join(dir, f.Name), os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			return err
		}
		if _, err := io.Copy(dst, src); err != nil {
			return err
		}
		dst.Close()
	}
	return nil
}
