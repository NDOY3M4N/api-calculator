### Usage

After pulling this docker image you can run it using this command

```sh
docker run \
    -p 3000:3000 \
    --env JWT_SECRET=my_secret_key \
    ghcr.io/NDOY3M4N/api-calculator:latest
```

> [!NOTE]
> If this command is too much for you 😉 you can just use the `compose.yml` provided in this repo.
