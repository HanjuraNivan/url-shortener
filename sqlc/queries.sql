-- name: GetShortURL :one
SELECT * FROM url_shortener.urlshortenermapping
WHERE shortURL = $1 LIMIT 1;

-- name: GetFullURL :one
SELECT * FROM url_shortener.urlshortenermapping
WHERE fullURL = $1 LIMIT 1;

-- name: CreateURLMapping :one
INSERT INTO url_shortener.urlshortenermapping (fullURL, shortURL, clickedTimes)
VALUES ($1, $2, 0)
RETURNING shortURL;

-- name: UpdateClickedTimes :exec
UPDATE url_shortener.urlshortenermapping
SET clickedTimes = clickedTimes + 1
WHERE shortURL = $1;

-- name: GetShortURLFromFullURL :one
SELECT shortURL FROM url_shortener.urlshortenermapping
WHERE fullURL = $1 LIMIT 1;

-- name: GetFullURLFromShortURL :one
SELECT fullURL FROM url_shortener.urlshortenermapping
WHERE shortURL = $1 LIMIT 1;    

-- name: GetClickedTimes :one
SELECT clickedTimes FROM url_shortener.urlshortenermapping      
WHERE shortURL = $1 LIMIT 1;