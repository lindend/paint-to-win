FROM debian
MAINTAINER Daniel Linden <d.linden88@gmail.com>
ADD wordlist wordlist

ENV ApiPort 80
ENV WordlistRoot /etc/paintToWin/wordlists

CMD ["sh", "-c", "./wordlist", "-db", "$DbConnectionString"]