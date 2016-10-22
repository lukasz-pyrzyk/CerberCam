using System.Drawing.Imaging;
using System.IO;
using System.Threading.Tasks;
using System.Windows.Input;
using CerberCam.Core;
using CerberCam.Core.Properties;

namespace CerberCam
{
    public class ViewModelBase
    {
        private readonly IQueueWrapper _wrapper = new QueueManager();

        private readonly bool _canExecute;
        private ICommand _clickCommand;

        public ViewModelBase()
        {
            _canExecute = true;
        }

        public ICommand SendCommand => _clickCommand ??
            (_clickCommand = new CommandHandler(() => Send(), _canExecute));

        public void Send()
        {
            byte[] data;
            using (MemoryStream ms = new MemoryStream())
            {
                Resources.golang_sh_600x600.Save(ms, ImageFormat.Jpeg);
                data = ms.ToArray();
            }

            Message msg = new Message
            {
                Email = "lukasz.pyrzyk@gmail.com",
                Photo = data
            };

            Task.Run(() => _wrapper.SendAsync(ref msg));
        }
    }
}