package main

const appName = "AppJar"
const appVersion = "0.1"

const baseDockerfile = `
FROM ubuntu:jammy

ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get -y update && apt-get -y install python3-uinput python3-pip git xvfb dbus-x11 libgstreamer1.0-dev libgstreamer-plugins-base1.0-dev libgstreamer-plugins-bad1.0-dev gstreamer1.0-plugins-base gstreamer1.0-plugins-good gstreamer1.0-plugins-bad gstreamer1.0-plugins-ugly gstreamer1.0-libav gstreamer1.0-tools gstreamer1.0-x gstreamer1.0-alsa gstreamer1.0-gl gstreamer1.0-gtk3 gstreamer1.0-qt5 gstreamer1.0-pulseaudio python3-pkgconfig libxxhash-dev libx11-dev libxfixes-dev libxext-dev libxdamage-dev libxkbfile-dev libxrandr-dev libxtst-dev libxcomposite-dev libxres-dev libgtk-3-dev python3-cairo-dev python-gi-dev liblz4-dev python3-dbus

RUN pip3 install --upgrade cython xdg pyxdg paramiko
RUN cd /tmp && git clone https://github.com/Xpra-org/xpra && cd xpra && git checkout v6.2.4 && python3 ./setup.py install
RUN cd /tmp && git clone https://github.com/Xpra-org/xpra-html5 && cd xpra-html5 && git checkout v17 && python3 ./setup.py install
WORKDIR /app
`

const startupCmdTemplate = `xpra start :100 --start-child "%s" --bind-tcp=0.0.0.0:%s --exit-with-children --auth=password --daemon=no --password=%s --mdns=no`
