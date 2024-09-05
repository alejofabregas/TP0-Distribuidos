import socket
import logging
import signal
import os
from common.utils import Bet, store_bets, load_bets, has_won

CLIENT_COUNT = int(os.getenv("CLIENT_COUNT"))

class Server:
    def __init__(self, port, listen_backlog):
        # Initialize server socket
        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._server_socket.bind(('', port))
        self._server_socket.listen(listen_backlog)
        self._clients = []
        self._finished_clients = {}

        # Set signal handlers
        self._shutdown_triggered = False
        signal.signal(signal.SIGTERM, self.__handle_shutdown)
        signal.signal(signal.SIGINT, self.__handle_shutdown)

    def run(self):
        """
        Dummy Server loop

        Server that accept a new connections and establishes a
        communication with a client. After client with communucation
        finishes, servers starts to accept new connections again
        """
        while not self._shutdown_triggered:
            try:
                client_sock = self.__accept_new_connection()
                self._clients.append(client_sock)
                self.__handle_client_connection(client_sock)
            except OSError as e:
                if not self._shutdown_triggered:
                    logging.error("action: accept_client | result: fail | error: {e}")
        
        logging.info('action: server_shutdown | result: success')

    def __handle_client_connection(self, client_sock):
        """
        Read message from a specific client socket and closes the socket

        If a problem arises in the communication with the client, the
        client socket will also be closed
        """
        try:
            length_bytes = self.__read_all(client_sock, 4) # Read 4 bytes (32 bits)
            length = int.from_bytes(length_bytes, "big")

            if length == 0:
                self.__handle_finish(client_sock)
                return

            msg = self.__read_all(client_sock, length).rstrip().decode('utf-8')
            addr = client_sock.getpeername()
            logging.info(f'action: receive_message | result: success | ip: {addr[0]}')

            try:
                bets = msg.split("\n")
                for i, bet in enumerate(bets):
                    bet_data = bet.split("|")
                    bet = Bet(*bet_data)
                    bets[i] = bet
                store_bets(bets)
                logging.info(f'action: apuesta_recibida | result: success | cantidad: {len(bets)}')
                response = "Batch OK\n".encode('utf-8')
                self.__write_all(client_sock, response)
                logging.info('action: ack_enviado | result: success')
            except:
                logging.error(f'action: apuesta_recibida | result: fail | cantidad: {len(bets)}')
                response = "Batch ERR\n".encode('utf-8')
                self.__write_all(client_sock, response)
                logging.info('action: error_enviado | result: success')

        except OSError as e:
            logging.error("action: receive_message | result: fail | error: {e}")
        finally:
            if length != 0:
                client_sock.close()

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

    def __handle_shutdown(self, signum, stack_frame):
        """
        Closes file descriptors and logs the server shutdown.
        Changes the value of the shutdown flag to True.
        This function is called when a SIGTERM or SIGINT signal is received.
        """
        logging.info(f'action: server_shutdown | result: in_progress | signal: {signum}')
        
        self._shutdown_triggered = True
        
        for client in self._clients:
            client.close()
            logging.info('action: client_shutdown | result: success')
        
        if self._server_socket:
            self._server_socket.close()

    def __read_all(self, client_sock, length_bytes):
        buffer = bytearray()
        while len(buffer) < length_bytes:
            partial_read = client_sock.recv(length_bytes - len(buffer))
            buffer.extend(partial_read)
        return bytes(buffer)

    def __write_all(self, client_sock, buffer):
        bytes_written = 0
        while bytes_written < len(buffer):
            n = client_sock.send(buffer[bytes_written:])
            bytes_written += n
        return bytes_written

    def __handle_finish(self, client_sock):
        if len(self._finished_clients) == CLIENT_COUNT - 1:
            logging.info('action: sorteo | result: success')
            bets = load_bets()
            winning_ids = []
            for bet in bets:
                if has_won(bet):
                    winning_ids.append(bet.document)
            msg = "|".join(winning_ids)
            logging.info(f'action: sorteo_ganadores | result: success | cant_ganadores: {len(winning_ids)}')
            
            # Add this last client
            addr = client_sock.getpeername()
            self._finished_clients[addr] = client_sock
            for client in self._finished_clients:
                self.__write_all(self._finished_clients[client], (msg + "\n").encode('utf-8'))
                self._finished_clients[client].close()
        else:
            addr = client_sock.getpeername()
            self._finished_clients[addr] = client_sock
            logging.info(f'action: new_finished_client | result: success | cant: {len(self._finished_clients)} | clients: {self._finished_clients}')
