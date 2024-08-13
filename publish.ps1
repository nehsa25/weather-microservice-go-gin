
$build = (get-date).ToString("yyyyMMddHHmm");

write-host "Building docker image for weather service ${build}";
docker build . -t nehsa/ascii-weather:$build --platform linux/amd64;

write-host "Pushing image to DockerHub...";
docker push nehsa/ascii-weather:$build;

write-host "Done!";