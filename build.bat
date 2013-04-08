rd bin /S /Q
cd server
goxc -os=linux -arch="amd64" -d="../bin" -z="false"
goxc -os=windows -arch="386" -d="../bin" -z="false"
cd ../client
goxc -os=linux -arch="amd64" -d="../bin" -z="false"
goxc -os=windows -arch="386" -d="../bin" -z="false"
cd ..
del "./bin/unknown/downloads.md"
mv "./bin/unknown/windows_386" "./bin"
mv "./bin/unknown/linux_amd64" "./bin"
rd "./bin/unknown" /S /Q
