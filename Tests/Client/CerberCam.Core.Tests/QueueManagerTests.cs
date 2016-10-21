using System.IO;
using System.Linq;
using System.Text;
using NSubstitute;
using ProtoBuf;
using RabbitMQ.Client;
using Xunit;

namespace CerberCam.Core.Tests
{
    public class QueueManagerTests
    {
        [Fact]
        public void SendAsync_SendsMessage()
        {
            // Arrange
            IModel channel = Substitute.For<IModel>();
            IConnection connection = Substitute.For<IConnection>();
            IConnectionFactory factory = Substitute.For<IConnectionFactory>();
            factory.CreateConnection().ReturnsForAnyArgs(connection);
            connection.CreateModel().ReturnsForAnyArgs(channel);
            var message = new Message { Email = "test@test.com", Photo = Encoding.UTF8.GetBytes("lorem ipsum") };
            byte[] expectedSerializationBytes;

            using (MemoryStream stream = new MemoryStream())
            {
                Serializer.Serialize(stream, message);
                expectedSerializationBytes = stream.ToArray();
            }

            IQueueWrapper wrapper = new QueueManager(factory);

            // Act
            wrapper.SendAsync(ref message);

            // Assert
            channel.Received(1).BasicPublish(Arg.Any<string>(), Arg.Any<string>(), Arg.Any<IBasicProperties>(), Arg.Is<byte[]>(x => x.SequenceEqual(expectedSerializationBytes)));
        }
    }
}
