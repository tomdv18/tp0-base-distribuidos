from .utils import Bet
import socket
import logging

def recieve_bet(client_sock):

    lenght_bytes = client_sock.recv(2)
    lenght = int.from_bytes(lenght_bytes, byteorder='big')
    logging.info(f'action: recieve_length | result: success | length: {lenght}')
    msg = client_sock.recv(lenght).decode('utf-8')
    addr = client_sock.getpeername()
    bets = msg.split('>')  
    

    bet_objects = []
    failed_bets = 0

    for bet in bets:
        if bet:  
            bet_details = bet.split(';')  
            if len(bet_details) != 6:
                logging.error(f'action: validate_bet | result: fail | invalid_bet: {bet}')
                failed_bets += 1
                continue
            

            bet_obj = Bet(*bet_details)
            bet_objects.append(bet_obj)
            logging.info(f'action: create_bet | result: success | bet: {bet_obj.number}')
    
    
    return bet_objects, failed_bets

def send_response(client_sock, bet_number):
    client_sock.send("{}\n".format(bet_number).encode('utf-8'))