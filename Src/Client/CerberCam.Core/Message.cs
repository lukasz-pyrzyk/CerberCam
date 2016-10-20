using ProtoBuf;

namespace CerberCam.Core
{
    [ProtoContract]
    public struct Message
    {
        [ProtoMember(1)]
        public string Email { get; set; }
    }
}
