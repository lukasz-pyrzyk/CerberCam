using System.Windows;

namespace CerberCam
{
    public partial class MainWindow : Window
    {
        public MainWindow()
        {
            InitializeComponent();
            DataContext = new ViewModelBase();
        }
    }
}
