$base = "http://localhost:8080/api"

function Req($label, $method, $url, $body=$null, $token=$null) {
    $h = @{}
    if ($token) { $h["Authorization"] = "Bearer $token" }
    try {
        $p = @{ Uri=$url; Method=$method; UseBasicParsing=$true; Headers=$h }
        if ($body) { $p.Body=$body; $p.ContentType="application/json" }
        $r = Invoke-WebRequest @p
        Write-Host "OK   $label -> $($r.StatusCode)"
        return $r.Content | ConvertFrom-Json
    } catch {
        $code = $_.Exception.Response.StatusCode.value__
        $msg  = $_.ErrorDetails.Message
        Write-Host "FAIL $label -> $code | $msg"
        return $null
    }
}

# Health
Req "GET /health" "GET" "$base/health" | Out-Null

# Signup
$s = Req "POST /auth/signup" "POST" "$base/v1/auth/signup" '{"email":"admin@example.com","password":"Password123","name":"Admin User"}'
$token = $s.data.accessToken
$refresh = $s.data.refreshToken

# Login
$l = Req "POST /auth/login" "POST" "$base/v1/auth/login" '{"email":"admin@example.com","password":"Password123"}'
if ($l) { $token = $l.data.accessToken; $refresh = $l.data.refreshToken }

# Auth protected
Req "GET  /auth/me" "GET" "$base/v1/auth/me" $null $token | Out-Null
Req "POST /auth/refresh" "POST" "$base/v1/auth/refresh" "{`"refreshToken`":`"$refresh`"}" | Out-Null
Req "POST /auth/forgot-password" "POST" "$base/v1/auth/forgot-password" '{"email":"admin@example.com"}' | Out-Null
Req "POST /auth/logout" "POST" "$base/v1/auth/logout" $null $token | Out-Null

# Re-login after logout
$l = Req "POST /auth/login (re)" "POST" "$base/v1/auth/login" '{"email":"admin@example.com","password":"Password123"}'
if ($l) { $token = $l.data.accessToken }

# App routes
Req "POST /app/auth/mentee-signup" "POST" "$base/v1/app/auth/mentee-signup" '{"email":"mentee@example.com","password":"Password123","name":"Mentee User"}' | Out-Null
Req "GET  /app/context" "GET" "$base/v1/app/context" $null $token | Out-Null
Req "GET  /app/sheets" "GET" "$base/v1/app/sheets" $null $token | Out-Null
Req "GET  /app/leaderboard" "GET" "$base/v1/app/leaderboard" $null $token | Out-Null
Req "GET  /app/me/profile" "GET" "$base/v1/app/me/profile" $null $token | Out-Null

# Org routes
$org = Req "POST /organizations" "POST" "$base/v1/organizations" '{"name":"Test Org","slug":"test-org","description":"A test org"}' $token
$orgId = $org.data.id
Req "GET  /organizations" "GET" "$base/v1/organizations" $null $token | Out-Null
Req "GET  /organizations/:id" "GET" "$base/v1/organizations/$orgId" $null $token | Out-Null
Req "PATCH /organizations/:id" "PATCH" "$base/v1/organizations/$orgId" '{"name":"Updated Org"}' $token | Out-Null

# Org members
Req "GET  /organizations/:id/members" "GET" "$base/v1/organizations/$orgId/members" $null $token | Out-Null

# Bootcamp routes
$bc = Req "POST /bootcamps" "POST" "$base/v1/organizations/$orgId/bootcamps" '{"name":"Test Bootcamp","description":"A test bootcamp"}' $token
$bootcampId = $bc.data.id
Req "GET  /bootcamps" "GET" "$base/v1/organizations/$orgId/bootcamps" $null $token | Out-Null
Req "GET  /bootcamps/:id" "GET" "$base/v1/organizations/$orgId/bootcamps/$bootcampId" $null $token | Out-Null
Req "PATCH /bootcamps/:id" "PATCH" "$base/v1/organizations/$orgId/bootcamps/$bootcampId" '{"name":"Updated Bootcamp"}' $token | Out-Null
Req "GET  /bootcamps/:id/enrollments" "GET" "$base/v1/organizations/$orgId/bootcamps/$bootcampId/enrollments" $null $token | Out-Null

# Problem routes
$pr = Req "POST /problems" "POST" "$base/v1/organizations/$orgId/problems" '{"title":"Two Sum","description":"Find two numbers","difficulty":"easy"}' $token
$problemId = $pr.data.id
Req "GET  /problems" "GET" "$base/v1/organizations/$orgId/problems" $null $token | Out-Null
Req "GET  /problems/:id" "GET" "$base/v1/organizations/$orgId/problems/$problemId" $null $token | Out-Null
Req "PATCH /problems/:id" "PATCH" "$base/v1/organizations/$orgId/problems/$problemId" '{"title":"Two Sum Updated"}' $token | Out-Null

# Tag routes
$tg = Req "POST /tags" "POST" "$base/v1/organizations/$orgId/tags" '{"name":"arrays"}' $token
$tagId = $tg.data.id
Req "GET  /tags" "GET" "$base/v1/organizations/$orgId/tags" $null $token | Out-Null

# Problem-Tag association
Req "POST /problems/:id/tags" "POST" "$base/v1/organizations/$orgId/problems/$problemId/tags" "{`"tagIds`":[`"$tagId`"]}" $token | Out-Null

# Resources
$rs = Req "POST /problems/:id/resources" "POST" "$base/v1/organizations/$orgId/problems/$problemId/resources" '{"title":"LeetCode","url":"https://leetcode.com"}' $token
$resourceId = $rs.data.id
Req "GET  /problems/:id/resources" "GET" "$base/v1/organizations/$orgId/problems/$problemId/resources" $null $token | Out-Null

# Assignment group
$ag = Req "POST /assignment-groups" "POST" "$base/v1/organizations/$orgId/bootcamps/$bootcampId/assignment-groups" '{"title":"Week 1","description":"First week"}' $token
$groupId = $ag.data.id
Req "GET  /assignment-groups" "GET" "$base/v1/organizations/$orgId/bootcamps/$bootcampId/assignment-groups" $null $token | Out-Null
Req "GET  /assignment-groups/:id" "GET" "$base/v1/organizations/$orgId/bootcamps/$bootcampId/assignment-groups/$groupId" $null $token | Out-Null
Req "POST /assignment-groups/:id/problems" "POST" "$base/v1/organizations/$orgId/bootcamps/$bootcampId/assignment-groups/$groupId/problems" "{`"problemIds`":[`"$problemId`"]}" $token | Out-Null

# Assignments
Req "GET  /assignments" "GET" "$base/v1/organizations/$orgId/bootcamps/$bootcampId/assignments" $null $token | Out-Null

# Doubts
Req "GET  /doubts" "GET" "$base/v1/doubts" $null $token | Out-Null
Req "GET  /doubts/me" "GET" "$base/v1/doubts/me" $null $token | Out-Null

# Analytics - leaderboard
Req "GET  /leaderboard" "GET" "$base/v1/bootcamps/$bootcampId/leaderboard" $null $token | Out-Null

# Analytics - polls
Req "GET  /polls" "GET" "$base/v1/bootcamps/$bootcampId/polls" $null $token | Out-Null

# Super admin routes
Req "GET  /super-admin/organizations" "GET" "$base/v1/super-admin/organizations" $null $token | Out-Null
Req "GET  /super-admin/bootcamps" "GET" "$base/v1/super-admin/bootcamps" $null $token | Out-Null
Req "GET  /super-admin/problems" "GET" "$base/v1/super-admin/problems" $null $token | Out-Null

Write-Host "`n--- IDs used ---"
Write-Host "orgId:       $orgId"
Write-Host "bootcampId:  $bootcampId"
Write-Host "problemId:   $problemId"
Write-Host "tagId:       $tagId"
Write-Host "groupId:     $groupId"
