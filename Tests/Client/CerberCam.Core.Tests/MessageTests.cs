using System.IO;
using System.Text;
using ProtoBuf;
using Xunit;

namespace CerberCam.Core.Tests
{
    public class MessageTests
    {
        [Fact]
        public void CanSerializeAndDeserialize()
        {
            // Arrange
            var message = new Message
            {
                Email = "test@test.com",
                Photo = Encoding.UTF8.GetBytes("lorem ipsum")
            };

            Message expected;
            // Act
            using (MemoryStream ms = new MemoryStream())
            {
                Serializer.Serialize(ms, message);
                ms.Seek(0, SeekOrigin.Begin);
                expected = Serializer.Deserialize<Message>(ms);
            }

            // Assert
            Assert.Equal(message.Photo, expected.Photo);
            Assert.Equal(message.Email, expected.Email);
        }

    }
}
