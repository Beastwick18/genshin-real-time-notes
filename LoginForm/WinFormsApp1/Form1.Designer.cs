namespace WebViewLogin
{
    partial class Form1
    {
        /// <summary>
        ///  Required designer variable.
        /// </summary>
        private System.ComponentModel.IContainer components = null;

        /// <summary>
        ///  Clean up any resources being used.
        /// </summary>
        /// <param name="disposing">true if managed resources should be disposed; otherwise, false.</param>
        protected override void Dispose(bool disposing)
        {
            if (disposing && (components != null))
            {
                components.Dispose();
            }
            base.Dispose(disposing);
        }

        #region Windows Form Designer generated code

        /// <summary>
        ///  Required method for Designer support - do not modify
        ///  the contents of this method with the code editor.
        /// </summary>
        private void InitializeComponent()
        {
            webView = new Microsoft.Web.WebView2.WinForms.WebView2();
            doneButton = new Button();
            helpLabel = new Label();
            refreshInput = new TextBox();
            refreshIntervalLabel = new Label();
            label1 = new Label();
            label2 = new Label();
            darkModeCheckBox = new CheckBox();
            ((System.ComponentModel.ISupportInitialize)webView).BeginInit();
            SuspendLayout();
            // 
            // webView
            // 
            webView.AllowExternalDrop = true;
            webView.CreationProperties = null;
            webView.DefaultBackgroundColor = Color.Black;
            webView.Dock = DockStyle.Fill;
            webView.Location = new Point(0, 0);
            webView.Name = "webView";
            webView.Size = new Size(354, 561);
            webView.TabIndex = 0;
            webView.ZoomFactor = 1D;
            webView.CoreWebView2InitializationCompleted += webview21_coreinit;
            webView.NavigationCompleted += webview21_NavigationCompleted;
            // 
            // doneButton
            // 
            doneButton.Anchor = AnchorStyles.Bottom | AnchorStyles.Left | AnchorStyles.Right;
            doneButton.BackColor = SystemColors.Control;
            doneButton.Location = new Point(177, 515);
            doneButton.Name = "doneButton";
            doneButton.Size = new Size(177, 46);
            doneButton.TabIndex = 1;
            doneButton.Text = "Done";
            doneButton.UseVisualStyleBackColor = false;
            doneButton.Click += doneButton_Click;
            // 
            // helpLabel
            // 
            helpLabel.Location = new Point(0, 0);
            helpLabel.Name = "helpLabel";
            helpLabel.Size = new Size(354, 42);
            helpLabel.TabIndex = 3;
            helpLabel.Text = "Once you have successfully logged in and entered the desired refresh interval, click \"Done\"";
            helpLabel.TextAlign = ContentAlignment.MiddleCenter;
            // 
            // refreshInput
            // 
            refreshInput.Location = new Point(119, 515);
            refreshInput.Name = "refreshInput";
            refreshInput.PlaceholderText = "Refresh Interval (sec)";
            refreshInput.Size = new Size(59, 23);
            refreshInput.TabIndex = 6;
            refreshInput.Text = "60";
            refreshInput.TextAlign = HorizontalAlignment.Center;
            // 
            // refreshIntervalLabel
            // 
            refreshIntervalLabel.Anchor = AnchorStyles.Bottom | AnchorStyles.Left;
            refreshIntervalLabel.Location = new Point(0, 515);
            refreshIntervalLabel.Name = "refreshIntervalLabel";
            refreshIntervalLabel.Size = new Size(121, 19);
            refreshIntervalLabel.TabIndex = 7;
            refreshIntervalLabel.Text = "Refresh Interval (sec):";
            refreshIntervalLabel.TextAlign = ContentAlignment.BottomCenter;
            // 
            // label1
            // 
            label1.Anchor = AnchorStyles.Bottom | AnchorStyles.Left;
            label1.Location = new Point(0, 534);
            label1.Name = "label1";
            label1.Size = new Size(121, 27);
            label1.TabIndex = 10;
            label1.Text = "Dark Mode:";
            label1.TextAlign = ContentAlignment.MiddleCenter;
            // 
            // label2
            // 
            label2.Anchor = AnchorStyles.Bottom | AnchorStyles.Left;
            label2.Location = new Point(119, 534);
            label2.Name = "label2";
            label2.Size = new Size(59, 27);
            label2.TabIndex = 11;
            label2.TextAlign = ContentAlignment.BottomCenter;
            // 
            // darkModeCheckBox
            // 
            darkModeCheckBox.AutoSize = true;
            darkModeCheckBox.Checked = true;
            darkModeCheckBox.CheckState = CheckState.Checked;
            darkModeCheckBox.Location = new Point(142, 541);
            darkModeCheckBox.Name = "darkModeCheckBox";
            darkModeCheckBox.Size = new Size(15, 14);
            darkModeCheckBox.TabIndex = 12;
            darkModeCheckBox.UseVisualStyleBackColor = true;
            // 
            // Form1
            // 
            AutoScaleDimensions = new SizeF(7F, 15F);
            AutoScaleMode = AutoScaleMode.Font;
            ClientSize = new Size(354, 561);
            Controls.Add(darkModeCheckBox);
            Controls.Add(label2);
            Controls.Add(label1);
            Controls.Add(refreshIntervalLabel);
            Controls.Add(refreshInput);
            Controls.Add(helpLabel);
            Controls.Add(doneButton);
            Controls.Add(webView);
            FormBorderStyle = FormBorderStyle.FixedSingle;
            MaximizeBox = false;
            MinimizeBox = false;
            Name = "Form1";
            Text = "Login to Hoyolab";
            Load += Form1_Load;
            ((System.ComponentModel.ISupportInitialize)webView).EndInit();
            ResumeLayout(false);
            PerformLayout();
        }

        #endregion

        private Microsoft.Web.WebView2.WinForms.WebView2 webView;
        private Button doneButton;
        private Label helpLabel;
        private Label refreshIntervalLabel;
        private Label label1;
        private TextBox refreshInput;
        private Label label2;
        private CheckBox darkModeCheckBox;
    }
}
