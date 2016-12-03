using System.ComponentModel;
using System.ComponentModel.DataAnnotations;
using System.Diagnostics;
using System.Drawing.Imaging;
using System.IO;
using System.Threading.Tasks;
using System.Windows.Input;
using CerberCam.Core;
using Emgu.CV;

namespace CerberCam
{
    public class ViewModelBase : INotifyPropertyChanged
    {
        private readonly IQueueWrapper _wrapper = new QueueManager();

        private readonly bool _canExecute;
        private ICommand _clickCommand;
        private string _email;

        private static readonly EmailAddressAttribute EmailValidator = new EmailAddressAttribute();

        public ViewModelBase()
        {
            _canExecute = true;
        }

        public ICommand SendCommand => _clickCommand ??
            (_clickCommand = new CommandHandler(() => Send(), _canExecute));

        public event PropertyChangedEventHandler PropertyChanged;
        
        public string Email
        {
            get { return _email; }
            set
            {
                if (value != _email)
                {
                    _email = value;
                    PropertyChanged?.Invoke(this, new PropertyChangedEventArgs(nameof(Email)));
                }
            }
        }

        public void Send()
        {
            if (string.IsNullOrEmpty(Email) || !EmailValidator.IsValid(Email))
            {
                Debug.WriteLine($"Invalid email: {Email}");
                return;
            }

            byte[] data;
            using (Capture c = new Capture())
            {
                c.Grab();
                Mat queryFrame = c.QueryFrame();

                using (var ms = new MemoryStream())
                {
                    queryFrame.Bitmap.Save(ms, ImageFormat.Jpeg);
                    data = ms.ToArray();
                }
            }

            Message msg = new Message
            {
                Email = Email,
                Photo = data
            };

            Task.Run(() => _wrapper.SendAsync(ref msg));
        }
    }
}