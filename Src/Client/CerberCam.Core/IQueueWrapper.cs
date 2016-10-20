namespace CerberCam.Core
{
    public interface IQueueWrapper
    {
        void SendAsync(Message msg);
    }
}
