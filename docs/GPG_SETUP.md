# GPG Setup for Signing Releases

This guide explains how to set up GPG signing for Terraform Provider releases.

## Why GPG Signing?

Terraform Registry **requires** signed releases for security. The GPG signature ensures:
- The release artifacts are authentic
- They haven't been tampered with
- They come from a verified source

## Prerequisites

- GPG installed: `gpg --version`
- GitHub account with admin access to the repository

## Step 1: Generate GPG Key

```bash
# Generate a new GPG key
gpg --full-generate-key

# Select:
# - Key type: (1) RSA and RSA
# - Key size: 4096
# - Expiration: 0 (does not expire) or set your preference
# - Real name: Your name or "SCINFRA Terraform Provider"
# - Email: your-email@example.com
# - Comment: "Code signing key for Terraform Provider"
```

## Step 2: List and Export Keys

```bash
# List your GPG keys
gpg --list-secret-keys --keyid-format=long

# Export the private key (replace KEY_ID with your key ID)
gpg --armor --export-secret-keys KEY_ID > private-key.asc

# Export the public key
gpg --armor --export KEY_ID > public-key.asc

# Get the fingerprint (needed for GitHub Actions)
gpg --fingerprint KEY_ID
```

## Step 3: Add GPG Keys to GitHub Secrets

Go to your repository settings:
`https://github.com/mrB4el/terraform-provider-aeza/settings/secrets/actions`

Add these secrets:

1. **GPG_PRIVATE_KEY**: 
   - Copy contents of `private-key.asc`
   - Include the `-----BEGIN PGP PRIVATE KEY BLOCK-----` and `-----END PGP PRIVATE KEY BLOCK-----` lines

2. **GPG_PASSPHRASE**: 
   - The passphrase you set when creating the key
   - If you didn't set one, leave empty

3. **GPG_FINGERPRINT** (optional, automatically detected):
   - The 40-character fingerprint without spaces
   - Example: `1234567890ABCDEF1234567890ABCDEF12345678`

## Step 4: Add the Public Key to Terraform Registry

```bash
gpg --armor --export KEY_ID
```

Add the exported public key under [Terraform Registry Signing Keys](https://registry.terraform.io/settings/gpg-keys) for your namespace. Adding it to GitHub is optional.

## Step 5: Publish to Keyserver (Optional but Recommended)

```bash
# Submit to Ubuntu keyserver
gpg --keyserver keyserver.ubuntu.com --send-keys KEY_ID

# Submit to OpenPGP keyserver
gpg --keyserver keys.openpgp.org --send-keys KEY_ID
```

## Step 6: Test Locally

```bash
# Test GoReleaser locally (without publishing)
goreleaser release --snapshot --clean

# Check the dist/ directory for build artifacts
ls -la dist/
```

## Step 7: Create a Test Release

```bash
# Create and push a test tag
git tag -a v0.1.1-test -m "Test release"
git push origin v0.1.1-test

# Check GitHub Actions: https://github.com/mrB4el/terraform-provider-aeza/actions
# Verify the release was created and signed
```

## Verification

Users can verify releases:

```bash
# Import the public key exported by the provider maintainer
gpg --import public-key.asc

# Verify the signature
gpg --verify terraform-provider-aeza_0.2.0_SHA256SUMS.sig \
             terraform-provider-aeza_0.2.0_SHA256SUMS

# Verify checksums
shasum -a 256 -c terraform-provider-aeza_0.2.0_SHA256SUMS
```

## Troubleshooting

### "gpg: signing failed: Inappropriate ioctl for device"

```bash
export GPG_TTY=$(tty)
```

### "No secret key" error

Make sure the private key is properly imported in GitHub Actions secrets.

### Key expiration

If your key expires, extend it:

```bash
gpg --edit-key KEY_ID
expire
# Set new expiration
save

# Re-export and update GitHub secrets
```

## Security Best Practices

1. **Store the private key securely** - Use a password manager
2. **Backup the key** - Store in a secure location
3. **Set expiration** - Consider 2-4 year expiration
4. **Revocation certificate** - Generate and store securely
5. **Use strong passphrase** - For the GPG key

## Resources

- [GitHub GPG signing](https://docs.github.com/en/authentication/managing-commit-signature-verification)
- [Terraform Registry Publishing](https://developer.hashicorp.com/terraform/registry/providers/publishing)
- [GoReleaser Documentation](https://goreleaser.com/customization/sign/)

