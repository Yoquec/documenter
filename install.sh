SHARE_DIR="/usr/share/documenter"

# Move assets
rm -rf $SHARE_DIR 2> /dev/null
mkdir -p $SHARE_DIR
cp -r static $SHARE_DIR/

# Move binary
rm -f /usr/bin/documenter 2> /dev/null
install -C -Dm 755 -t "/usr/bin" src/documenter
