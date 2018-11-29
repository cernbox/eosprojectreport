# 
# eosprojectreport spec file 
#

Name: eosprojectreport
Summary: Utility to get information for CERNBox project spaces
Version: 0.0.1
Release: 1%{?dist}
License: MIT
BuildRoot: %{_tmppath}/%{name}-buildroot
Group: CERN-IT/ST
BuildArch: x86_64
Source: %{name}-%{version}.tar.gz

%description
This RPM provides the eosprojectreport tool, an utility to obtain project space
information without wasting your time.

# Don't do any post-install weirdness, especially compiling .py files
%define __os_install_post %{nil}

%prep
%setup -n %{name}-%{version}

%install
# server versioning

# installation
rm -rf %buildroot/
mkdir -p %buildroot/usr/local/bin
install -m 755 eosprojectreport	     %buildroot/usr/local/bin/eosprojectreport

%clean
rm -rf %buildroot/

%preun

%post

%files
%defattr(-,root,root,-)
/usr/local/bin/*

%changelog
* Thu Nov 29 2018 Hugo Gonzalez Labrador <hugo.gonzalez.labrador@cern.ch> 0.0.1
- First release
