# Invalid email.
curl http://localhost:8080/appmetas \
    --include \
    --header "Content-Type: application/yaml" \
    --request "POST" \
    --data 'title: App w/ Invalid maintainer email
version: 1.0.1
maintainers:
- name: Firstname Lastname
  email: apptwohotmail.com
company: Upbound Inc.
website: https://upbound.io
source: https://github.com/upbound/repo
license: Apache-2.0
description: |
 ### blob of markdown
 More markdown'

# Missing fields.
curl http://localhost:8080/appmetas \
    --include \
    --header "Content-Type: application/yaml" \
    --request "POST" \
    --data 'title: App w/ Invalid maintainer email
version: 1.0.1'

# Query all.
curl http://localhost:8080/appmetas:query \
    --include \
    --header "Content-Type: application/yaml" \
    --request "POST" \
    --data '{}'

# Query one.
curl http://localhost:8080/appmetas:query \
    --include \
    --header "Content-Type: application/yaml" \
    --request "POST" \
    --data 'title: Valid App 1
version: 0.0.1'

# Valid ones.
curl http://localhost:8080/appmetas \
    --include \
    --header "Content-Type: application/yaml" \
    --request "POST" \
    --data 'title: Valid App 1
version: 0.0.1
maintainers:
- name: firstmaintainer app1
  email: firstmaintainer@hotmail.com
- name: secondmaintainer app1
  email: secondmaintainer@gmail.com
company: Random Inc.
website: https://website.com
source: https://github.com/random/repo
license: Apache-2.0
description: |
 ### Interesting Title
 Some application content, and description'

curl http://localhost:8080/appmetas \
    --include \
    --header "Content-Type: application/yaml" \
    --request "POST" \
    --data 'title: Valid App 2
version: 1.0.1
maintainers:
- name: AppTwo Maintainer
  email: apptwo@hotmail.com
company: Upbound Inc.
website: https://upbound.io
source: https://github.com/upbound/repo
license: Apache-2.0
description: |
 ### Why app 2 is the best
 Because it simply is...'
