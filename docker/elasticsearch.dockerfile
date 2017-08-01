FROM elasticsearch:5.5.1

ENTRYPOINT ["/docker-entrypoint.sh"]
CMD ["elasticsearch"]