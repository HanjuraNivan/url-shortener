CREATE schema if not exists url_shortener;

CREATE TABLE if not exists url_shortener.urlshortenermapping (
  id uuid primary key default gen_random_uuid() not null,
  created_at TIMESTAMPTZ default current_timestamp not null,
  fullURL varchar unique not null,
  shortURL varchar unique not null,
  clickedTimes integer
);