# creds
```creds``` is a tool for managing the creds.txt file in penetration tests.


```
go install github.com/r0nk/creds@latest
```
## get
Get the path of creds.txt file, which can be in the current directory or any parent directory.
```
creds path

echo "bob:hunter1" | anew $(creds path)
```

print all usernames
```
creds usernames
```

print usernames for matching password
```
creds usernames 'hunter1'
```

print all passwords
```
creds passwords
```

print password for username bob
```
creds passwords bob
sshpass $(creds passwords bob | head -n 1) ssh bob@hostname echo SUCCESS!
for pass in $(creds passwords bob) ; do sshpass $pass ssh bob@hostname ls; done
```


## generate
Output creds with first letter of username words capitalized
(bob dude:pass) -> (Bob Dude:pass)
```
creds title
```

OUTPUT CREDENTIALS BUT WITH COOL KID CRUISE CONTROL
```
creds capslock
```

output creds in hip modern teenage engineering design.
```
creds lowercase
```

Output creds with every password for each username
```
creds permutate
```

Output creds with username as password (admin:admin)
```
creds dual
for c in $(creds dual); do wget -r ftp://$c@127.0.0.1; done
```

Output creds with all possible permutations
```
creds all
```
