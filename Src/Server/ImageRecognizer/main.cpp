//
// Created by bentoo on 24.11.16.
//

#include <tensorflow/c_api.h>
#include <cstdio>
#include <fstream>
#include <cstring>
#include <cassert>

static const char* model_file = "Data/tensorflow_inception_graph.pb";
static const char* label_file = "Data/imagenet_comp_graph_label_strings.txt";

struct File
{
	FILE*			handle;
	std::size_t		length;
};

void checkStatus(TF_Status* status)
{
	auto code = TF_GetCode(status);
	if(code != TF_OK)
	{
		printf("Code %d : %s", code, TF_Message(status));
		assert(false && "Stoped on Tensorflow error!");
	}
}

File openFile(const char* path)
{
	File file;
	file.handle = fopen(path, "r");
	file.length = 0;

	if(file.handle != nullptr)
	{
		file.length = (std::size_t) feof(file.handle);
	}

	return file;
}

int main()
{
	TF_Status* status = TF_NewStatus();
	TF_Graph* graph = TF_NewGraph();

	File model = openFile(model_file);

	TF_Buffer* buffer = TF_NewBufferFromString(model.handle, model.length);

	TF_ImportGraphDefOptions* importOptions = TF_NewImportGraphDefOptions();

	TF_GraphImportGraphDef(graph, buffer, importOptions, status);

	checkStatus(status);

	TF_Session* session = TF_NewSession(TF_NewSessionOptions(), status);

	checkStatus(status);

	return 0;
}
