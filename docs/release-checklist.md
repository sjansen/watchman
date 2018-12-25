# Release Checklist

1) Verify that all tests are passing.
1) Update `CHANGELOG.md` and commit.
1) Create release branch:

    ```sh
    git checkout -b release/v0.1
    ```

1) Tag release.

    ```sh
    git tag -a v0.1.0 -m "Release 0.1.0"
    ```

1) Push commits and tags.

    ```sh
    git push origin release/v0.1
    git push origin v0.1.0
    ```

1) Update GitHub.

    * https://github.com/sjansen/watchman/releases
