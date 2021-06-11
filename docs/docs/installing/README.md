# Installing Drago

We are working to make Drago available as a pre-compiled binary or as a package for several operating systems. 

Meanwhile, you can build Drago from source.

## Compiling from Source

To compile from source, you will need Go installed and configured properly.

1. Clone the Drago repository:
    ```bash
    $ git clone https://github.com/seashell/drago.git
    $ cd drago
    ```
2. Download all necessary Go modules into the module cache:
   
    ```bash
    $ go mod download
    ```
3. Build Drago for your current system.

    ```bash
    $ go build -o ./build/drago
    ```

## Verify the binary

To verify Drago was built correctly, run the binary:

```bash
$ ./build/drago
```
