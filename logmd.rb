class Logmd < Formula
  desc "A minimal, local-first journal CLI written in Go"
  homepage "https://github.com/hellodizzy/logmd"
  version "1.0.0"
  
  # Update these URLs when you create GitHub releases
  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/hellodizzy/logmd/releases/download/v1.0.0/logmd-darwin-arm64"
      sha256 "5f91b8e7bbdeec78645af88776af3c65b4b83ce1e623b617e1b704a4fbd44e26"
    else
      url "https://github.com/hellodizzy/logmd/releases/download/v1.0.0/logmd-darwin-amd64"
      sha256 "78f843d1504bcf0769a5f6f2f5be631efc432a031d509ae653ace36566850329"
    end
  end

  on_linux do
    if Hardware::CPU.arm?
      url "https://github.com/hellodizzy/logmd/releases/download/v1.0.0/logmd-linux-arm64"
      sha256 "bd2371b8646b0642416faac096937868f209f6b9f7a4aeb12763c41625ec7528"
    else
      url "https://github.com/hellodizzy/logmd/releases/download/v1.0.0/logmd-linux-amd64"
      sha256 "afb5338ce826dcdfea4ac87a152665998209d0ee9b65937341cde26b713e0b9d"
    end
  end

  def install
    bin.install Dir["*"].first => "logmd"
  end

  test do
    assert_match "logmd", shell_output("#{bin}/logmd --version")
    assert_match "Display current configuration", shell_output("#{bin}/logmd config --help")
  end
end
