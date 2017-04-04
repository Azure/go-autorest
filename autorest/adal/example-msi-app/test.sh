#!/usr/bin/env bash

# RUN THIS ON THE REMOTE MACHINE WITH MSI EXTENSION ACTIVE

wget 'https://colemickmsitest0.blob.core.windows.net/testapp/example-msi-app'
chmod +x example-msi-app
./example-msi-app
