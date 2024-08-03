from django import forms
from django.core.validators import EmailValidator
from django_recaptcha.fields import ReCaptchaField
from django_recaptcha.widgets import ReCaptchaV2Checkbox

class ContactForm(forms.Form):
    name = forms.CharField(max_length=50, widget=forms.TextInput(attrs={'class': 'form-control'}))
    email = forms.EmailField(widget=forms.EmailInput(attrs={'class': 'form-control'}), validators=[EmailValidator()])
    message = forms.CharField(widget=forms.Textarea(attrs={'class': 'form-control'}))
    captcha = ReCaptchaField(widget=ReCaptchaV2Checkbox) 
