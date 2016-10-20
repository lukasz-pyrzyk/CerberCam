using System.Threading.Tasks;
using System.Windows.Input;
using CerberCam.Core;

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
            Message msg = new Message { Email = "lukasz.pyrzyk@gmail.com" };
            _wrapper.SendAsync(msg);
        }
    }
}