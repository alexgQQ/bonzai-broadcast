SERVICE_NAME=foobar
EXECUTABLE=foobar

rsync_dirs=""

install_service="
    cp ./$SERVICE_NAME /lib/systemd/system  \n
    cp ./$EXECUTABLE /usr/bin               \n
    systemctl daemon-reload                 \n
    systemctl enable $SERVICE_NAME          \n
    systemctl start $SERVICE_NAME           \n
"

uninstall_service="
    systemctl stop $SERVICE_NAME            \n
    systemctl disable $SERVICE_NAME         \n
    rm /lib/systemd/system/$SERVICE_NAME    \n
    rm /usr/bin/$EXECUTABLE                 \n
    systemctl daemon-reload                 \n
"

echo $install_service
