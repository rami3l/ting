class Ting < Formula
    desc "Yet another tcping."
    homepage "https://github.com/rami3l/ting"
    version "{version}"
    url "{url_mac}"
    sha256 "{sha256_mac}"

    if OS.linux?
      url "{url_linux}"
      sha256 "{sha256_linux}"
    end

    def install
      bin.install "ting"
    end
end