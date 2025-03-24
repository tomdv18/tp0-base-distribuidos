from .utils import Bet
import socket
import logging

def recieve_bet(client_sock):

    lenght_bytes = client_sock.recv(2)
    lenght = int.from_bytes(lenght_bytes, byteorder='big')
    logging.info(f'action: recieve_length | result: success | length: {lenght}')
    msg = client_sock.recv(lenght).decode('utf-8')
    addr = client_sock.getpeername()
    logging.info(f'action: receive_message | result: success | ip: {addr[0]} | msg: {msg}')
    bet = Bet(*msg.split(';'))
    logging.info(f'action: create_bet | result: success | bet: {bet.number}')
    
    return bet

def send_response(client_sock, bet_number):
    client_sock.send("{}\n".format(bet_number).encode('utf-8'))