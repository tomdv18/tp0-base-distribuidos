from .utils import Bet
import socket
import logging

NEW_BET = 0x00
WINNERS = 0x01


def recieve_message(client_sock):

    message_type = client_sock.recv(1)
    if not message_type:
        logging.error("not message_type")
        return False, None,0
    
    if message_type == NEW_BET:
        return recieve_bet(client_sock)
    elif message_type == WINNERS:
        return recieve_winners(client_sock)
        
    else:
        logging.error(f"Unknown message type {message_type}")
        return False, None, 0

def recieve_winners(client_sock):
    lenght_bytes = client_sock.recv(2)
    lenght = int.from_bytes(lenght_bytes, byteorder='big')
    logging.info(f'action: recieve_winners | result: success | length: {lenght}')
    msg = client_sock.recv(lenght).decode('utf-8')

    return True, None, msg



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
    

    return False, bet_objects, failed_bets

def send_response(client_sock, msg):
    client_sock.send("{}\n".format(msg).encode('utf-8'))

def send_winners_response(winners, client_sock):
    