using ProtoBuf;

namespace CerberCam.Core
{
    [ProtoContract]
    public struct Message
    {
        [ProtoMember(1)]
        public string Email { get; set; }

        [ProtoMember(2)]
        public byte[] Photo { get; set; }
    }
}
