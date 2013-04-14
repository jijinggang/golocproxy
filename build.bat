cd server
goxc -os=linux -arch="amd64" -d="../bin" -z="false"
goxc -os=windows -arch="386" -d="../bin" -z="false"
cd ../client
goxc -os=linux -arch="amd64" -d="../bin" -z="false"
goxc -os=windows -arch="386" -d="../bin" -z="false"
cd ../bin
rd windows_386 /S /Q
mv "./unknown/windows_386" "."
rd linux_amd64 /S /Q
mv "./unknown/linux_amd64" "."
rd "./unknown" /S /Q
cd ..