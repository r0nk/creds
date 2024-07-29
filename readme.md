```creds``` is a tool for managing the creds.txt file in penetration tests.

# get
path of creds.txt file
```
creds path
```

print all usernames
```
creds usernames
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

print usernames for matching password hunter1
```
creds usernames 'hunter1'
```

# mod
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
```

Output creds with all possible permutations
```
creds all
```

mod flags

```
 -i     write changes directly to the creds file
 -u     only apply to username
 -p     only apply to password

```
