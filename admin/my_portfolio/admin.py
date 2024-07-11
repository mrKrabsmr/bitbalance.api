from django.contrib import admin
from django.contrib.auth.models import Group

from my_portfolio.models import User, Portfolio

admin.site.site_header = "Мой портфель"
admin.site.index_title = "Администрирование сайта"

admin.site.unregister(Group)

class PortfolioInline(admin.TabularInline):
    model = Portfolio
    extra = 0

@admin.register(User)
class UserAdmin(admin.ModelAdmin):
    list_display = ("username",)
    inlines = (PortfolioInline,)


@admin.register(Portfolio)
class PortfolioAdmin(admin.ModelAdmin):
    list_display = ("user", "cryptocurrency", "cryptocurrency_symbol")

    def get_queryset(self, request):
        queryset = super().get_queryset(request)
        return queryset.select_related("user")
