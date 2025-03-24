import socket
import logging
import sys
import signal
from .utils import Bet, store_bets
from .comunications import recieve_bet, send_response


class Server:
    def __init__(self, port, listen_backlog):
        # Initialize server socket
        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._server_socket.bind(('', port))
        self._server_socket.listen(listen_backlog)

        signal.signal(signal.SIGINT, self.__signal_handler)
        signal.signal(signal.SIGTERM, self.__signal_handler)

        self.clientes = []

    def run(self):
        """
        Dummy Server loop

        Server that accept a new connections and establishes a
        communication with a client. After client with communucation
        finishes, servers starts to accept new connections again
        """

        # TODO: Modify this program to handle signal to graceful shutdown
        # the server
        while True:
            client_sock = self.__accept_new_connection()
            self.clientes.append(client_sock)
            self.__handle_client_connection(client_sock)

    def __handle_client_connection(self, client_sock):
        """
        Read message from a specific client socket and closes the socket

        If a problem arises in the communication with the client, the
        client socket will also be closed
        """
        try:

            bet = recieve_bet(client_sock)
            
            store_bets([bet])

            logging.info(f'action: apuesta_almacenada | result: success | dni: {bet.document} |numero: {bet.number}')

            send_response(client_sock, bet.number)
        except OSError as e:
            logging.error("action: receive_message | result: fail | error: {e}")
        finally:
            client_sock.close()
            self.clientes.remove(client_sock)

    def __accept_new_connection(self):
        """
        Accept new connections

        Function blocks until a connection to a client is made.
        Then connection created is printed and returned
        """

        # Connection arrived
        logging.info('action: accept_connections | result: in_progress')
        c, addr = self._server_socket.accept()
        logging.info(f'action: accept_connections | result: success | ip: {addr[0]}')
        return c

    def __signal_handler(self, signum, frame):
        logging.info(f"action: signal_received | result: success | signal: {signum}")
        
        self._server_socket.close()
        logging.debug("action: server socket closed | result: success")

        # Cerrar todos los sockets de clientes
        for cliente in self.clientes:
            cliente.close()

        logging.debug("action: close all client sockets | result: success")
        sys.exit(0)