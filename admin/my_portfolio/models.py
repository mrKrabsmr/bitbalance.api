import uuid

from django.db import models

from django.contrib.auth.hashers import make_password
from django.contrib.auth.models import BaseUserManager
from django.contrib.auth.base_user import AbstractBaseUser


class UserManager(BaseUserManager):

    def create_user(self, username, password):
        if not username:
            raise ValueError("User must have username")

        user = self.model(username=username, password=password)
        user.save()

        return user

    def create_superuser(self, username, password, is_superuser=True):
        if not username:
            raise ValueError("User must have username")

        user = self.model(
            username=username, password=password, is_superuser=is_superuser, is_staff=True,
        )
        user.save()

        return user


class User(AbstractBaseUser):
    id = models.UUIDField(primary_key=True, default=uuid.uuid4, editable=False)
    username = models.CharField(unique=True, max_length=255)
    password = models.CharField(max_length=255)
    created_at = models.DateTimeField(auto_now_add=True)
    is_superuser = models.BooleanField(default=False)
    is_staff = models.BooleanField(default=False)

    last_login = None

    USERNAME_FIELD = "username"

    REQUIRED_FIELDS = ("password",)

    objects = UserManager()

    class Meta:
        managed = False
        db_table = 'users'
        verbose_name = "пользователь"
        verbose_name_plural = "пользователи"

    def __str__(self):
        return self.username

    def save(self, *args, **kwargs):
        if self._state.adding or self.password != self.__class__.objects.get(pk=self.pk).password:
            self.password = make_password(self.password)
        return super().save(*args, **kwargs)

    def has_perm(self, perm, obj=None):
        return self.is_superuser

    def has_module_perms(self, app_label):
        return self.is_superuser


class Portfolio(models.Model):
    id = models.UUIDField(primary_key=True, default=uuid.uuid4, editable=False)
    user = models.ForeignKey('User', models.CASCADE)
    cmc_cryptocurrency_id = models.IntegerField()
    cryptocurrency = models.CharField(max_length=255, null=True, blank=True)
    cryptocurrency_symbol = models.CharField(
        max_length=30, null=True, blank=True)
    price = models.FloatField()
    count = models.FloatField()
    purchase_time = models.DateTimeField()
    commentary = models.CharField(blank=True, null=True)
    created_at = models.DateTimeField(auto_now_add=True)

    class Meta:
        managed = False
        db_table = 'portfolios'
        verbose_name = "портфель"
        verbose_name_plural = "портфель"
