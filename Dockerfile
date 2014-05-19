from base/arch
run pacman -Suy --noconfirm
run pacman -S --noconfirm wget
run wget -qO- https://storage.googleapis.com/golang/go1.2.2.linux-amd64.tar.gz | tar vxzC /usr/local
env GOROOT /usr/local/go
env PATH /usr/local/go/bin:$PATH
run pacman -S --noconfirm mercurial git llvm gcc
run mkdir /usr/local/gopath
env GOPATH "/usr/local/gopath"
RUN go get -v github.com/goskydome/mqtt-stack
RUN cd $GOPATH/src/github.com/goskydome/mqtt-stack/
RUN go get -d -v ./...
RUN go build -v ./...
RUN go install
