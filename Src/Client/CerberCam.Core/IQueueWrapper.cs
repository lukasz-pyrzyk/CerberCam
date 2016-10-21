namespace CerberCam.Core
{
    public interface IQueueWrapper
    {
        void SendAsync(ref Message msg);
    }
}
