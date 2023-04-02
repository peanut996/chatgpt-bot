import logging

from engine.route import app, set_session
from engine.session.session import Session
from engine.tool import load_config


def main():
    logging.basicConfig(level=logging.INFO)
    app.logger.setLevel(logging.WARNING)
    config = load_config()
    session = Session(config=config)
    port = config['engine']['port']
    debug = config['engine'].get('debug', False)
    set_session(session)
    app.run(host="0.0.0.0", port=port, debug=debug)


if __name__ == "__main__":
    main()
