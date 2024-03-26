from django.contrib.auth.hashers import BCryptPasswordHasher

class MPBCryptPasswordHasher(BCryptPasswordHasher):
    rounds = 8 
