# spots_go
A Golang project to test out web scraping and download videos using golang and youtube-dl

# How to use:
Simply put in the command:
go run main.go PLAYLIST_URL TARGETFOLDER

and the program will start to download the music versions of the songs.

You can chain multiple urls:

go run main.go PLAYLIST_URL_1 PLAYLIST_URL_2 ... PLAYLIST_URL_N TARGET_FOLDER

## Todo:
1. Replace goquery with Selenium to scrap JS modifed html
2. Add better error handling 
