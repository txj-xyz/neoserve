# File Server

Provided by Golang with Auth, publically serves the `public/` folder as a web server.


# Examples

## ShareX Custom Uploader

Below is a ShareX custom uploader, replace the `{{}}` values respectively `TODO: replace with templating CI/etc`
```json
{
  "Version": "15.0.0",
  "Name": "Neoserve File Host",
  "DestinationType": "ImageUploader",
  "RequestMethod": "POST",
  "RequestURL": "https://{{ config.server.domain }}/v1/upload",
  "Headers": {
    "Authorization": "Bearer {{ config.auth.token }}"
  },
  "Body": "MultipartFormData",
  "FileFormName": "file",
  "URL": "{json:url}"
}
```


## Uploads

```shell
# Curl request
curl --location 'http://localhost:8081/v1/upload' \
--header 'Authorization: Bearer {{ config.auth.token }}' \
--form 'file=@"/C:/Users/timmyboy/Pictures/omg_a_meme.jpg"'
```


## File Listing

`Requires: config.server.devMode = true`


```shell
curl --location 'http://localhost:8081/v1/files'
```