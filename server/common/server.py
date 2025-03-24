import socket
import logging
import sys
import signal
import threading
from .utils import Bet, store_bets
from .comunications import recieve_message, send_response


class Server:
    def __init__(self, port, listen_backlog):
        # Initialize server socket
        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._server_socket.bind(('', port))
        self._server_socket.listen(listen_backlog)

        signal.signal(signal.SIGINT, self.__signal_handler)
        signal.signal(signal.SIGTERM, self.__signal_handler)

        self.sockets_clientes = []
        self.clientes_finalizados = []
        self.clientes_fin_lock = threading.Lock()

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
            self.sockets_clientes.append(client_sock)
            self.__handle_client_connection(client_sock)

    def __handle_client_connection(self, client_sock):
        """
        Read message from a specific client socket and closes the socket

        If a problem arises in the communication with the client, the
        client socket will also be closed
        """
        try:

            is_winners, bets, failed = recieve_message(client_sock)

            if not is_winners:
                for bet in bets:
                    store_bets([bet])

                msg =""
                if failed > 0:
                    msg =f'action: apuesta_recibida | result: fail | cantidad: {len(bets) + failed}'
                else:
                    msg =f'action: apuesta_recibida | result: success | cantidad: {len(bets)}'

                logging.info(msg)

                send_response(client_sock, msg)
            else:
                with self.clientes_fin_lock:
                    self.clientes_finalizados.append(client_sock)
                    logging.info("action: finalizar_cliente | result: success")
                
                if len(self.clientes_finalizados) == len(self.sockets_clientes):
                    logging.info("action: get winner | result: in_progress")
                    
                pass
        except OSError as e:
            logging.error("action: receive_message | result: fail | error: {e}")
        finally:
            client_sock.close()
            self.sockets_clientes.remove(client_sock)

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
        for cliente in self.sockets_clientes:
            cliente.close()

        logging.debug("action: close all client sockets | result: success")
        sys.exit(0)