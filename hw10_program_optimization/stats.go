package hw10programoptimization

import (
	"bufio"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"io"
	"regexp"
	"strings"
)

type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	u, err := getUsers(r)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}
	return countDomains(u, domain)
}

type users [100_000]User

func getUsers(r io.Reader) (users, error) {
	scanner := bufio.NewScanner(r)
	var result users

	var user User
	i := 0
	for scanner.Scan() {
		if err := json.Unmarshal(scanner.Bytes(), &user); err != nil {
			return result, err
		}
		result[i] = user

		user = User{}
		i++
	}

	return result, nil
}

func countDomains(u users, domain string) (DomainStat, error) {
	result := make(DomainStat)
	re, err := regexp.Compile("\\." + domain)
	if err != nil {
		return nil, err
	}

	for _, user := range u {
		matched := re.MatchString(user.Email)

		if matched {
			domain := strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])
			result[domain]++
		}
	}
	return result, nil
}
