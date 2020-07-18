# video-streamer
HLS video streamer developed in  GO with Fiber (https://github.com/gofiber/fiber)

Convert mp4 files to TLS compatible format using following command and place it in specific folder in public folder.

`ffmpeg -i movie.mp4 -profile:v baseline -level 3.0 -s 640x360 -start_number 0 -hls_time 10 -hls_list_size 0 -f hls index.m3u8`
