using System.Threading.Tasks;

namespace CerberCam.Core
{
    public interface IQueueWrapper
    {
        Task SendAsync(Message msg);
    }
}
