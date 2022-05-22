# s3-cisbench

![image](https://user-images.githubusercontent.com/7364201/161427647-8b1978de-a37c-4a57-a3df-fde00e8e5737.png)

A simple command line tool that checks S3 bucket against (security-) best
practices, mainly CIS benchmark based.

> ⚠ ⚠  This project is currently work-in-progress

## CIS AWS Benchmark v1.4.0: Storage

The AWS Benchmark section 'Storage' contains the S3 bucket related items,
namely:

* 2.1.1 Ensure all S3 buckets employ encryption-at-rest
* 2.1.2 Ensure S3 Bucket Policy is set to deny HTTP requests
* 2.1.3 Ensure MFA Delete is enable on S3 buckets
* 2.1.4 Ensure all data in Amazon S3 has been discovered (_out of scope_)
* 2.1.5 Ensure that S3 Buckets are configured with 'Block public access'
  * ✖ ✔ BlockPublicAcls (BPA)
  * ✖ ✔ BlockPublicPolicy (BPP)
  * ✖ ✔ IgnorePublicAcls (IPA)
  * ✖ ✔ RestrictPublicBuckets (RPB)

Currently known limitations:
* encryption at rest only checks for default AES256 algorithm and reports false otherwise

## Usage

```text
$ s3-cisbench --help
s3-csibench is a tool that analyses S3 bucket against CIS benchmark rules.

Usage:
  s3-cisbench [flags]
  s3-cisbench [command]

Available Commands:
  audit       Audit S3 buckets against applicable CIS benchmark items
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  list        List AWS S3 buckets.

Flags:
  -d, --debug   Enable verbose logging
  -h, --help    help for s3-cisbench

Use "s3-cisbench [command] --help" for more information about a command.
```

### 'audit' Command Output Example

<img src="docs/s3-cisbench.png" width="60%">

### 'audit -o json' Command Example with 'jq' processing

Usage of json output with leveraging `jq` for further filtering:

<img src="docs/s3-cisbench-json.png"   width="90%">

## Build

```sh
git clone https://github.com/rollwagen/s3-cisbench
cd s3-cisbench
make
```
