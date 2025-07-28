#!/bin/bash
#
#
# gets the original owner of a pod
#

set -x

getowner(){
  local object=$1
  local owner=$(kubectl get $object -o jsonpath='{ .metadata.ownerReferences[].kind}/{.metadata.ownerReferences[].name}' || return 1)
 

  if [ -z $owner ]
  then
    return 1
  fi

  if [ $owner == '/' ] 
  then 
    echo $object
    return 0
  else
    getowner $owner
    return $?
  fi
}
kubectl get event -w --field-selector reason=Started,involvedObject.kind=Pod \
       	-o go-template='{{printf "%s\n" .involvedObject.name}}' |
(
while true
do
  read line 
  if [ ! -z $line ]
  then
    owner=$(getowner pod/$line)
    if [[ ! -z $owner && ! $owner == / ]]
    then 
      sleep 10
      if  kubectl delete pod $line
      then
	  sleep 10
          kubectl get $owner -o yaml |grep -v -E 'image:.*-baffs$' | sed 's/image:.*$/&-baffs/' | kubectl apply -f -
      fi
    fi
  fi
done
)
