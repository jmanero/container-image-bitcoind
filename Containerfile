ARG VERSION="27.1"

FROM docker.io/library/fedora:latest AS fetch
ARG VERSION

## Download and validate a bitcoin-core build
RUN curl -O https://bitcoincore.org/bin/bitcoin-core-${VERSION}/bitcoin-${VERSION}-$(arch)-linux-gnu.tar.gz
ADD https://bitcoincore.org/bin/bitcoin-core-${VERSION}/SHA256SUMS .
ADD https://bitcoincore.org/bin/bitcoin-core-${VERSION}/SHA256SUMS.asc .

ADD https://github.com/bitcoin-core/guix.sigs/archive/refs/heads/main.tar.gz ./signers.tar.gz
RUN mkdir /signers; tar -xzf signers.tar.gz --strip-components 1 -C /signers

## gpg forks some agent daemons that leave lock files laying around, which block subsequent gpg calls
RUN gpg --lock-never --import /signers/builder-keys/* && gpgconf --kill all
RUN gpg --lock-never --verify SHA256SUMS.asc SHA256SUMS && gpgconf --kill all
RUN sha256sum --ignore-missing --check SHA256SUMS

## Build a skelton for the scratch output image
RUN mkdir -p /build/etc /build/usr/bin /build/usr/lib /build/usr/lib64
RUN ln -s /usr/bin /build/bin
RUN ln -s /usr/lib /build/lib
RUN ln -s /usr/lib64 /build/lib64

## Add bitcoind and bitcoin-cli to the output image
RUN tar -xvzf bitcoin-${VERSION}-$(arch)-linux-gnu.tar.gz --strip-components 1 -C /build/usr bitcoin-${VERSION}/bin/bitcoin{d,-cli}

## Add the dynamic loader and library dependencies to the image
# $ ldd ./build/bin/bitcoind
# linux-vdso.so.1 (0x00007ffc9234a000)
# libpthread.so.0 => /lib64/libpthread.so.0 (0x00007081df64b000)
# libm.so.6 => /lib64/libm.so.6 (0x00007081df568000)
# libgcc_s.so.1 => /lib64/libgcc_s.so.1 (0x00007081df53b000)
# libc.so.6 => /lib64/libc.so.6 (0x00007081df34e000)
# /lib64/ld-linux-x86-64.so.2 (0x00007081e058b000)

RUN cp -aLv /usr/lib64/libpthread.so.0 /usr/lib64/libm.so.6 /usr/lib64/libgcc_s.so.1 /usr/lib64/libc.so.6 /build/usr/lib64

## Location of the dynamic loader varies across architectures.
RUN cp /usr/lib64/ld-linux* /build/usr/lib64 && ln -s /usr/lib64/ld-linux-* build/usr/bin/ld.so || true
RUN cp /usr/lib/ld-linux* /build/usr/lib && ln -s /usr/lib/ld-linux-* build/usr/bin/ld.so || true
RUN [ -f /build/usr/bin/ld.so ] || (echo "Unable to find a dynamic loader library in /usr/lib or /usr/lib64" && exit 1)

RUN cp -av /etc/ld.so.* /build/etc

## Build a scratch output image
FROM scratch

COPY --from=fetch /build /
COPY bitcoin.conf /etc/

ENTRYPOINT [ "/usr/bin/bitcoind" ]
CMD [ "-conf=/etc/bitcoin.conf" ]
