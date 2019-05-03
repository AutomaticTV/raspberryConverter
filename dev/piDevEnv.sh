wget https://storage.googleapis.com/golang/go1.12.4.linux-armv6l.tar.gz
sudo tar -C /usr/local -xvf go1.12.4.linux-armv6l.tar.gz
rm go1.12.4.linux-armv6l.tar.gz
cat >> ~/.bashrc << 'EOF'
export GOPATH=$HOME/go
export PATH=/usr/local/go/bin:$PATH:$GOPATH/bin
EOF
source ~/.bashrc
sudo apt-get install libpcre3 fonts-freefont-ttf omxplayer fbi -y
