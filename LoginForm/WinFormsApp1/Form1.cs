using System.Data;
using System.Text.Json;
using System.Text.Json.Serialization;
using Microsoft.Web.WebView2.Core;
using Microsoft.Web.WebView2.WinForms;

namespace WebViewLogin
{
    public partial class Form1 : Form
    {
        public Form1()
        {
            InitializeComponent();

            string[] args = Environment.GetCommandLineArgs();
            string url = "https://www.google.com";
            App = "unknown";
            if (args.Length > 1)
            {
                switch (args[1])
                {
                    case "genshin":
                        url = "https://act.hoyolab.com/app/community-game-records-sea/index.html#/ys";
                        Icon = Properties.Resources.GenshinIcon;
                        UidSelectorJS = "document.querySelector('.uid').innerHTML.split(' ')[0].substr(3)";
                        ExpectedUidLength = 9;
                        break;
                    case "hsr":
                        url = "https://act.hoyolab.com/app/community-game-records-sea/rpg/index.html#/hsr";
                        Icon = Properties.Resources.HsrIcon;
                        UidSelectorJS = "document.querySelector('.uid').innerHTML.split(' ')[0].substr(3)";
                        ExpectedUidLength = 9;
                        break;
                    case "zzz":
                        url = "https://act.hoyolab.com/app/mihoyo-zzz-game-record/index.html#/zzz";
                        Icon = Properties.Resources.ZzzIcon;
                        UidSelectorJS = "document.querySelector(\"[class^=uid_]\").innerHTML.trim().split('\\n')[0].substr(3)";
                        ExpectedUidLength = 10;
                        break;
                }
                App = args[1];
            }
            webView.Source = new Uri(url, UriKind.Absolute);
        }
        private string App { get; set; }

        private uint ExpectedUidLength { get; set; }

        private string UidSelectorJS { get; set; }

        private void Form1_Load(object sender, EventArgs e)
        {
            if (Screen.PrimaryScreen == null)
            {
                return;
            }
            var w = Screen.PrimaryScreen.WorkingArea.Width;
            var h = Screen.PrimaryScreen.WorkingArea.Height;
            var tw = Width;
            var th = Height;
            Location = new Point(w - tw, h - th);
        }

        private void Warning(string msg, string title = "")
        {
            Task.Run(() =>
            {
                MessageBox.Show(msg, title, MessageBoxButtons.OK, MessageBoxIcon.Warning);
            });
        }
        private async Task<Dictionary<string, CoreWebView2Cookie>> GetCookiesAsDict()
        {
            var cookies = await webView.CoreWebView2.CookieManager.GetCookiesAsync("https://act.hoyolab.com");
            return cookies.GroupBy(x => x.Name)
                          .ToDictionary(x => x.Key, x => x.First());
        }

        private async void doneButton_Click(object sender, EventArgs e)
        {
            var dict = await GetCookiesAsDict();

            if (!dict.TryGetValue("ltoken_v2", out var ltoken) || !dict.TryGetValue("ltuid_v2", out var ltuid))
            {
                Warning("Could not find ltoken_v2 or ltuid_v2 cookies. Make sure you are logged in.");
                return;
            }

            // string queryString = String.Format("document.querySelector('{}').innerHTML.split(' ')[0].substr(3)", this.UidSelector);
            string uid = await webView.ExecuteScriptAsync(UidSelectorJS);
            if (uid == "null" || uid == null)
            {
                Warning("Could not find UID. Make sure you are logged in.");
                return;
            }
            if (uid.Length < 6)
            {
                Warning("Could not find a valid UID: Got \"" + uid + "\"");
                return;
            }

            switch (this.App)
            {
                case "genshin":
                    uid = uid.Substring(1, uid.Length - 2);
                    break;
                case "hsr":
                    uid = uid.Substring(1, uid.Length - 2);
                    break;
                case "zzz":
                    uid = uid.Substring(3, uid.Length - 4);
                    break;
                default:
                    break;
            };

            if (!uid.All(char.IsAsciiDigit))
            {
                Warning("UID must be only numbers: Got \"" + uid + "\"");
                return;
            }

            if (uid.Length != ExpectedUidLength)
            {
                Warning("UID is not the expected length (" + ExpectedUidLength + "): Got \"" + uid + "\"");
                return;
            }

            uint refresh = 60;
            if (!uint.TryParse(refreshInput.Text, out refresh) || refresh < 30)
            {
                Warning("Not a valid refresh interval, must enter a positive integer greater than 30.");
                return;
            }



            var cookies = new Cookies()
            {
                RefreshInterval = refresh,
                UID = uid,
                Ltoken = ltoken.Value,
                Ltuid = ltuid.Value,
                DarkMode = darkModeCheckBox.Checked,
            };

            string json = JsonSerializer.Serialize(cookies, new JsonSerializerOptions()
            {
                WriteIndented = true
            });
            File.WriteAllText(App + "_cookie.json", json);
            Application.Exit();
        }

        private void webview21_coreinit(object sender, CoreWebView2InitializationCompletedEventArgs e)
        {
            webView.CoreWebView2.CookieManager.DeleteAllCookies();
            webView.CoreWebView2.Settings.UserAgent = "Mozilla/5.0 (Linux; Android 11; SAMSUNG SM-G973U) AppleWebKit/537.36 (KHTML, like Gecko) SamsungBrowser/14.2 Chrome/87.0.4280.141 Mobile Safari/537.36";
        }

        private void webview21_NavigationCompleted(object sender, CoreWebView2NavigationCompletedEventArgs e)
        {
            // Hide Scrollbars
            if (e.IsSuccess)
            {
                ((WebView2)sender).ExecuteScriptAsync("document.querySelector('body').style.overflow='scroll';var style=document.createElement('style');style.type='text/css';style.innerHTML='::-webkit-scrollbar{display:none}';document.getElementsByTagName('body')[0].appendChild(style)");
                // ((WebView2)sender).ExecuteScriptAsync("document.querySelector('.cp-y-no_login-btn').click()");
            }
        }
    }

    internal class Cookies
    {
        [JsonPropertyName("uid")]
        public required string UID { get; set; }

        [JsonPropertyName("refresh_interval")]
        public required uint RefreshInterval { get; set; }

        [JsonPropertyName("ltoken_v2")]
        public required string Ltoken { get; set; }

        [JsonPropertyName("ltuid_v2")]
        public required string Ltuid { get; set; }

        [JsonPropertyName("dark_mode")]
        public required Boolean DarkMode { get; set; }
    }
}
